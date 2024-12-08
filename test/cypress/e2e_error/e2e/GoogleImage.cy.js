describe('Google Image Search', () => {
  it('should be searching image on google', () => {
    //Navigate to google search page
    cy.visit('https://www.google.com/');
    //Adding wait time to load the page
    cy.wait(500);
    //Clicking on Images link on the google search page
    cy.contains('Images').click();
    //Entering search value in the textbox and performing th eenter key simulation
    cy.get('.gLFyf').type('Tesgrid automation{enter}');
  });
});
