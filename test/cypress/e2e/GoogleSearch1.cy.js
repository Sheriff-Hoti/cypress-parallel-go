describe('Google Search 1', () => {
  it('should be searching keyword on google', () => {
    //Navigate to google search page
    cy.visit('https://www.google.com/');
    //Adding wait time to load the page
    cy.wait(500);
    //Entering search value in the textbox and performing th eenter key simulation
    cy.get('.gLFyf').type('TestGrid{enter}');
  });
});
